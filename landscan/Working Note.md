Working notes on the landscan package
====================================

31 october 2018
----------------

**Getting the sample file**
eastview has provided a sample haiti file from 2014's landscan

The sample comes as directory with 2 subdirectories
- Country
- Population

No technical documentation is provided
The "Country" directory contains a "dbf" file and 2 subdirectories
- countries
- info
"countries" contains some "adf" files and a "metadat.xml" files
"info" contains one "arc.dir" file and 3 pairs of "arc000x.dat" and "arc000x.nit" files

The "Population" directory contains a "lsop2014.lyr" file and 2 directories
- info
- lspop2014

"info" contains one "arc.dir" file and 3 pairs of "arc000x.dat" and "arc000x.nit" files
"lsop2014" contains some "adf" files and a "metadat.xml" files

All this is cryptic for me. Let's try to understand. Notice that no ".shp" is present.

**understanding ESRI files**
reading http://www.esri.com/library/whitepapers/pdfs/shapefile.pdf

"A shapefile stores nontopological geometry and attribute information for the spatial features in a data set"
That is chinese to me. "nontopological geometry" ? "spatial features" ? "data set" ?

"The geometry for a feature is stored as a shape comprising a set of vector coordinates." again.

"An ESRI shapefile consists of a main file, an index file, and a dBASE table." That I understand.

"The main file is a direct access, variable-record-length file in which each record describes a shape with a list of its vertices.". Arghh, "direct access" ? "variable record length" ? "record" ? "shape" ? "vertice" ?

"In the index file, each record contains the offset of the corresponding main file record from the beginning of the main file." Arghh, I should learn what is a "record".

"The dBASE table contains feature attributes with one record per feature." A feature, according to wikipeadia is a type of geometric shape "a polygone" for instance.

"Attribute records in the dBASE file must be in the same order as records in the main file."

- Main file: counties.shp
- Index file: counties.shx
- dBASE table:   counties.dbf

Many files in the haiti landscan are not of those types